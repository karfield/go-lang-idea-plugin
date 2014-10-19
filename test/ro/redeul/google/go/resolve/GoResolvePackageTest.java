package ro.redeul.google.go.resolve;

import com.intellij.psi.PsiElement;
import ro.redeul.google.go.lang.psi.toplevel.GoTypeSpec;

import static ro.redeul.google.go.util.GoPsiTestUtils.getAs;

public class GoResolvePackageTest extends GoPsiResolveTestCase {

    @Override
    protected String getTestDataRelativePath() {
        return super.getTestDataRelativePath() + "package/";
    }

//    public void testBuiltinTypes() throws Exception {
//        doTest();
//    }
//
//    public void testBuiltinConversion() throws Exception {
//        doTest();
//    }
//
//    public void testVarBuiltinType() throws Exception {
//        doTest();
//    }
//
//    public void testVarMethodType() throws Exception {
//        doTest();
//    }
//
    public void testNamedImport() throws Exception {
        doTest();
    }
}
