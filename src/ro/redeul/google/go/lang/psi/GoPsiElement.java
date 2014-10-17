package ro.redeul.google.go.lang.psi;

import com.intellij.codeInsight.lookup.LookupElementBuilder;
import com.intellij.psi.PsiElement;
import ro.redeul.google.go.lang.psi.visitors.GoElementVisitor;
import ro.redeul.google.go.lang.psi.visitors.GoElementVisitorWithData;

/**
 * Author: Toader Mihai Claudiu <mtoader@gmail.com>
 * <p/>
 * Date: Jul 24, 2010
 * Time: 10:24:11 PM
 */
public interface GoPsiElement extends PsiElement {

    GoPsiElement[] EMPTY_ARRAY = new GoPsiElement[0];

    void accept(GoElementVisitor visitor);

    <T> T accept(GoElementVisitorWithData<T> visitor);

    void acceptChildren(GoElementVisitor visitor);

    LookupElementBuilder getCompletionPresentation();

    LookupElementBuilder getCompletionPresentation(GoPsiElement child);

    String getPresentationText();

    String getPresentationTailText();

    String getPresentationTypeText();

    GoPsiElement getReferenceContext();
}

