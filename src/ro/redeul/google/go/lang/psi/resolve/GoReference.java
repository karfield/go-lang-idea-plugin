package ro.redeul.google.go.lang.psi.resolve;

import com.intellij.psi.PsiElement;
import com.intellij.psi.PsiReference;
import org.jetbrains.annotations.NotNull;
import ro.redeul.google.go.lang.psi.GoPsiElement;

public abstract class GoReference<
        E extends GoPsiElement,
        Self extends GoReference<E, Self>> implements PsiReference {

    E element;

    @Override
    public PsiElement getElement() {
        return element;
    }

    @Override
    public boolean isSoft() {
        return false;
    }

    protected abstract Self self();

    @NotNull
    @Override
    public String getCanonicalText() {
        return getElement().getText();
    }

}
